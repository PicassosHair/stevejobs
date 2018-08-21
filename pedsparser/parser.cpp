//
//  parser.hpp
//  pedsparser
//
//  Created by Hao Liu on 8/16/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//

#include <stdlib.h>
#include <stdio.h>
#include "parser.hpp"
#include "database.hpp"
#include <iostream>
#include <fstream>
#include <nlohmann/json.hpp>
#include "mysql_connection.h"

#include <cppconn/driver.h>
#include <cppconn/exception.h>
#include <cppconn/resultset.h>
#include <cppconn/statement.h>

using namespace std;
using json = nlohmann::json;

namespace idsguard {
    
    vector<string> split(const string& str, const string& delim) {
        vector<string> tokens;
        size_t prev = 0, pos = 0;
        do
        {
            pos = str.find(delim, prev);
            if (pos == string::npos) pos = str.length();
            string token = str.substr(prev, pos-prev);
            if (!token.empty()) tokens.push_back(token);
            prev = pos + delim.length();
        }
        while (pos < str.length() && prev < str.length());
        return tokens;
    }
    
    Parser::Parser(const string& in_file_path, Database* db) {
        in_file_path_ = in_file_path;
        db_ = db;
    }
    
    Parser::~Parser() {}
    
    int Parser::StartWork() {
        ifstream inFile;
        inFile.open(in_file_path_);
        
        if (!inFile) {
            cout << "Unable to open file."
            << '\n';
            return -1;
        }
        
        std::size_t count = 0;
        
        auto on_parse = [&count, this](int depth, json::parse_event_t event, json& parsed) {
            if (event == json::parse_event_t::object_end) {
                if (parsed.find("patentRecord") != parsed.end()) {
                    string record = parsed["patentRecord"][0].dump();
                    ++count;
                    int status = this->updateApplications(record);
                    if (status != 0) {
                        cout << "Update Applications failed.\n";
                    }
                    return false;
                }
            }
            return true;
        };
        
        json::parse(inFile, on_parse);
        
        std::cout << "Found " << count << " applications.\n";
        
        inFile.close();
        return 0;
    }
    
    int Parser::updateApplications(const string& record) {
        json j3;
        try {
            j3 = json::parse(record);
        } catch (json::exception &e) {
            cout << "# ERR: " << e.what();
            return -1;
        }
        const string& applId = j3["patentCaseMetadata"]["applicationNumberText"]["value"];
        const json& peds_data = j3["patentCaseMetadata"];
        json& transactions = j3["prosecutionHistoryDataOrPatentTermData"];
        const json& invention_title_content = peds_data["inventionTitle"]["content"];
        if (!invention_title_content.is_array() ||
            invention_title_content.size() == 0 ||
            invention_title_content.is_null()) {
            cout << "Can't find the title of the application." << '\n';
            return -1;
        }
        
        const string& title = invention_title_content[0];

        string application_id;
        int status = db_->GetOrCreateApplication(applId, peds_data.dump(),
                                                 title, &application_id);
        if (status != 0) return -1;
        
        status = updateTransactions(application_id, transactions);
        
        return 0;
    }
    
    int Parser::updateTransactions(const string& application_id,
                                   json& transactions) {
        if (!transactions.is_array()) {
            return -1;
        }
        for (json::iterator it = transactions.begin();
             it != transactions.end(); ++it) {
            const auto& transaction = it.value();
            
            if (!transaction.is_object()) {
                return -1;
            }
            const string& recorded_date = transaction["recordedDate"];
            const string& desc_text = transaction["caseActionDescriptionText"];
            const string& delim = " , ";
            vector<string> desc;
            desc = split(desc_text, delim);
            if (desc.size() != 2) {
                cout << "Transaction description cannot be split to two.\n";
                cout << "Descrition size: " << desc.size() << '\n';
                cout << "Raw description: " << desc_text << '\n';
                return -1;
            }
            string description = desc[0];
            string code = desc[1];
            
            int status = db_->CreateOrIgnoreTransaction(application_id,
                                                        code,
                                                        description,
                                                        recorded_date);
            if (status != 0) {
                cout << "Failed to update Transaction.\n";
            }
        }
        return 0;
    }
}


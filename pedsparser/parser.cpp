//
//  parser.hpp
//  pedsparser
//
//  Created by Hao Liu on 8/16/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//

#ifndef parser_hpp
#define parser_hpp

#include <stdlib.h>
#include <stdio.h>
#include "parser.hpp"
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
    
    Parser::Parser(string& in_file_path) {
        in_file_path_ = in_file_path;
    }
    
    Parser::~Parser() {}
    
    void Parser::StartWork() {
        ifstream inFile;
        inFile.open(in_file_path_);
        
        if (!inFile) {
            cout << "Unable to open file."
            << '\n';
            exit(1);
        }
        
        char c;
        int left_curly_ct = 0;
        bool read_record = false;
        string record;
        string start_word;
        
        while (inFile >> c) {
            start_word += c;
            if (c == '"') {
                if (start_word == "\"patentRecord\"") {
                    read_record = true;
                }
                start_word = "\"";
            }
            
            if (c == '{' && read_record) {
                left_curly_ct++;
                
            } else if (c == '}' && read_record) {
                left_curly_ct--;
                
                if (left_curly_ct == 0) {
                    record += c;
                    auto j3 = json::parse(record);
                    cout << j3["patentCaseMetadata"]["applicationNumberText"]["value"] << '\n' << '\n';
                    record.clear();
                }
            }
            
            if (left_curly_ct > 0 && read_record) {
                record += c;
            }
        }
        
        inFile.close();
    }
    
    void Parser::ConnectDB() {
        cout << "Parser connecting database." << '\n';
        try {
            sql::Driver *driver;
            sql::Connection *con;
            sql::Statement *stmt;
            sql::ResultSet *res;
            
            /* Create a connection */
            driver = get_driver_instance();
            con = driver->connect("tcp://127.0.0.1:3306", "root", "");
            /* Connect to the MySQL test database */
            con->setSchema("idsguard_dev");
            
            stmt = con->createStatement();
            res = stmt->executeQuery("SELECT 'Hello World!' AS _message");
            while (res->next()) {
                cout << "\t... MySQL replies: ";
                /* Access column data by alias or column name */
                cout << res->getString("_message") << endl;
                cout << "\t... MySQL says it again: ";
                /* Access column data by numeric offset, 1 is the first column */
                cout << res->getString(1) << endl;
            }
            delete res;
            delete stmt;
            delete con;
            
        } catch (sql::SQLException &e) {
            cout << "# ERR: SQLException in " << __FILE__;
            cout << "(" << __FUNCTION__ << ") on line "
            << __LINE__ << endl;
            cout << "# ERR: " << e.what();
            cout << " (MySQL error code: " << e.getErrorCode();
            cout << ", SQLState: " << e.getSQLState() << " )" << endl;
        }
    }
}


#endif /* parser_hpp */

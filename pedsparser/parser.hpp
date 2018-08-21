//
//  parser.cpp
//  pedsparser
//
//  Created by Hao Liu on 8/16/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//
#ifndef parser_hpp
#define parser_hpp

#include "mysql_connection.h"
#include <string>
#include <nlohmann/json.hpp>
#include "database.hpp"

namespace idsguard {
    class Parser {
    private:
        std::string in_file_path_;
        Database* db_;
        int updateApplications(const std::string& record);
        int updateTransactions(const std::string& application_id,
                               nlohmann::json& transactions);
    
    public:
        Parser(const std::string& in_file_path, Database* database);
        ~Parser();
        int StartWork();
    };
}

#endif /* parser_hpp */

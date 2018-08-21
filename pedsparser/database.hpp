//
//  database.hpp
//  pedsparser
//
//  Created by Hao Liu on 8/19/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//
#ifndef database_hpp
#define database_hpp

#include <stdlib.h>
#include <stdio.h>
#include <string>
#include <cppconn/driver.h>
#include <cppconn/exception.h>
#include <cppconn/resultset.h>
#include <cppconn/statement.h>

#include "mysql_connection.h"

namespace idsguard {
    class Database {
    private:
        sql::Driver* driver_;
        sql::Connection* connection_;
        sql::ResultSet* executeSQLQuery(std::ostringstream& raw_query);
        bool executeActionSQLQuery(std::ostringstream& raw_query);
    public:
        Database(const std::string& host_name, const std::string& username,
                 const std::string& password, const std::string& database_name);
        ~Database();
        int GetOrCreateApplication(const std::string& applId,
                                   const std::string& peds_data,
                                   const std::string& title,
                                   std::string* application_id);
        int GetOrCreateTransactionCode(const std::string& code,
                                       const std::string& desc,
                                       std::string* transaction_code_id);
        int CreateOrIgnoreTransaction(const std::string& application_id,
                                      const std::string& code,
                                      const std::string& desc,
                                      const std::string& recorded_date);
    };
}

#endif /* database_hpp */

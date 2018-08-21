//
//  database.cpp
//  pedsparser
//
//  Created by Hao Liu on 8/19/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//
#include "database.hpp"
#include "mysql_connection.h"
#include <stdio.h>
#include <boost/algorithm/string.hpp>

#include <cppconn/driver.h>
#include <cppconn/exception.h>
#include <cppconn/resultset.h>
#include <cppconn/statement.h>

using namespace std;

namespace idsguard {
    // Clear query string stream.
    inline void clearQuery(std::ostringstream& sql) {
        sql.seekp(0);
        sql.clear();
        sql.str("");
    }
    
    Database::Database(const string& host_name, const string& username,
                       const string& password, const string& database_name) {
        driver_ = get_driver_instance();
        connection_ = driver_->connect(host_name, username, password);
        connection_->setSchema(database_name);
    }
    Database::~Database() {
        delete connection_;
    }
    
    sql::ResultSet* Database::executeSQLQuery(std::ostringstream& raw_query) {
        sql::ResultSet* result = nullptr;
        try {
            sql::Statement* statement;
            statement = connection_->createStatement();
            
            result = statement->executeQuery(raw_query.str().c_str());
            clearQuery(raw_query);
            
            delete statement;
        } catch (sql::SQLException &e) {
            cout << "# ERR: Incorrect SQL query: "
            << raw_query.str().c_str() << '\n';
            cout << "# ERR: SQLException in " << __FILE__;
            cout << "(" << __FUNCTION__ << ") on line "
            << __LINE__ << endl;
            cout << "# ERR: " << e.what();
            cout << " (MySQL error code: " << e.getErrorCode();
            cout << ", SQLState: " << e.getSQLState() << " )" << endl;
            
            clearQuery(raw_query);
            return nullptr;
        }
        return result;
    }
    
    bool Database::executeActionSQLQuery(std::ostringstream& raw_query) {
        bool result = false;
        try {
            sql::Statement* statement;
            statement = connection_->createStatement();
            
            result = statement->execute(raw_query.str().c_str());
            clearQuery(raw_query);
            
            delete statement;
        } catch (sql::SQLException &e) {
            cout << "# ERR: Incorrect SQL query: "
            << raw_query.str().c_str() << '\n';
            cout << "# ERR: SQLException in " << __FILE__;
            cout << "(" << __FUNCTION__ << ") on line "
            << __LINE__ << endl;
            cout << "# ERR: " << e.what();
            cout << " (MySQL error code: " << e.getErrorCode();
            cout << ", SQLState: " << e.getSQLState() << " )" << endl;
            
            clearQuery(raw_query);
            return false;
        }
        return true;
    }
    
    int Database::GetOrCreateApplication(const std::string& applId,
                                         const std::string& peds_data,
                                         const std::string& title,
                                         std::string* application_id) {
        sql::ResultSet* check_application_exist_result = nullptr;
        sql::ResultSet* get_application_result = nullptr;
        bool update_result, create_result;
        std::ostringstream query;
        query << "SELECT `id` AS _id FROM `Applications` WHERE applId = '"
        << applId << "' LIMIT 1;";
        check_application_exist_result = executeSQLQuery(query);
        if (check_application_exist_result == nullptr) return -1;
        
        if (check_application_exist_result->first()) {
            *application_id = check_application_exist_result->getString("_id");
            // Application exists. Update it.
            query
            << "UPDATE `Applications` SET `updatedAt` = NOW(), `pedsData` = '"
            << peds_data
            << "', `title` = '"
            << title << "' "
            << "WHERE applId = '" << applId << "';";
            update_result = executeActionSQLQuery(query);
        } else {
            // Create a new application.
            query
            << "INSERT INTO `Applications` (`createdAt`, `updatedAt`, `applId`, `pedsData`, `title`) "
            << "VALUES(NOW(), NOW(), '"
            << applId << "', '" << peds_data << "', '" << title << "'"
            << ");";
            create_result = executeActionSQLQuery(query);
            
            // Retrieve the application ID which is just created.
            query << "SELECT LAST_INSERT_ID() AS _id;";
            get_application_result = executeSQLQuery(query);
            if (get_application_result == nullptr) return -1;
            
            if (get_application_result->next()) {
                *application_id = get_application_result->getString("_id");
            } else {
                return -1;
            }
        }
        
        delete check_application_exist_result;
        delete get_application_result;
        return 0;
    }
    
    int Database::GetOrCreateTransactionCode(const std::string &code,
                                             const std::string &desc,
                                             std::string *transaction_code_id) {
        sql::ResultSet* check_code_result = nullptr;
        sql::ResultSet* get_code_result = nullptr;
        bool insert_result;
        std::ostringstream query;
        
        // Check if transaction code exists.
        query << "SELECT id AS _id FROM `TransactionCodes` WHERE `code` = '"
        << code << "';";
        check_code_result = executeSQLQuery(query);
        if (check_code_result == nullptr) return -1;
        
        if (check_code_result->next()) {
            *transaction_code_id = check_code_result->getString("_id");
        } else {
            query
            << "INSERT INTO `TransactionCodes` (`createdAt`, `updatedAt`, `code`, `description`, `type`, `initiator`, `isActionable`) "
            << "VALUES(NOW(), NOW(), '"
            << code << "', \"" << desc << "\", 'info', 'uspto', TRUE);";
            insert_result = executeActionSQLQuery(query);
            
            query << "SELECT LAST_INSERT_ID() AS _id;";
            
            get_code_result = executeSQLQuery(query);
            if (get_code_result == nullptr) return -1;
            
            if (check_code_result->next()) {
                *transaction_code_id = check_code_result->getString("_id");
            } else {
                return -1;
            }
        }
        delete check_code_result;
        delete get_code_result;
        return 0;
    }
    
    int Database::CreateOrIgnoreTransaction(const std::string &application_id,
                                            const std::string &code,
                                            const std::string &desc,
                                            const std::string &recorded_date) {
        sql::ResultSet* check_transaction_result = nullptr;
        sql::ResultSet* create_result = nullptr;
        std::string code_id;
        std::ostringstream query;
        
        int status = GetOrCreateTransactionCode(code, desc, &code_id);
        if (status != 0) return -1;
        
        // Check if transaction exists.
        query << "SELECT id FROM `Transactions` WHERE `transactionCodeId` = "
        << code_id
        << " AND `applicationId` = "
        << application_id
        << " AND DATE(`recordDate`) = '" << recorded_date << "';";
        check_transaction_result = executeSQLQuery(query);
        if (check_transaction_result == nullptr) return -1;
        
        if (!check_transaction_result->next()) {
            // Create a new transaction.
            query
            << "INSERT INTO `Transactions` (`createdAt`, `updatedAt`, `transactionCodeId`, `applicationId`, `recordDate`) "
            << "VALUES (NOW(), NOW(), "
            << code_id << ", "
            << application_id << ", '" << recorded_date << "');";
            create_result = executeSQLQuery(query);
            if (create_result == nullptr) return -1;
        }
        
        delete check_transaction_result;
        delete create_result;
        return 0;
    }
}

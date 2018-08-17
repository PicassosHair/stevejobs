//
//  parser.hpp
//  pedsparser
//
//  Created by Hao Liu on 8/16/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//

#ifndef parser_hpp
#define parser_hpp

#include <stdio.h>
#include "parser.hpp"
#include <iostream>
#include <fstream>

namespace idsguard {
    Parser::Parser(std::string& in_file_path, std::string& out_file_path) {
        in_file_path_ = in_file_path;
        out_file_path_ = out_file_path;
    }
    
    Parser::~Parser() {}
    
    void Parser::StartWork() {
        std::ifstream inFile;
        std::ofstream outFile (out_file_path_);
        inFile.open(in_file_path_);
        
        if (!inFile) {
            std::cout << "Unable to open file."
            << '\n';
            exit(1);
        }
        
        char c;
        int left_curly_ct = 0;
        bool read_record = false;
        std::string record;
        std::string start_word;
        
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
                    outFile << record << '\n';
                    record.clear();
                }
            }
            
            if (left_curly_ct > 0 && read_record) {
                record += c;
            }
        }
        
        outFile.close();
        inFile.close();
    }
}


#endif /* parser_hpp */

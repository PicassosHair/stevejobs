//
//  parser.cpp
//  pedsparser
//
//  Created by Hao Liu on 8/16/18.
//  Copyright © 2018 Hao Liu. All rights reserved.
//

#include <string>

namespace idsguard {
    class Parser {
        private:
            std::string in_file_path_;
            std::string out_file_path_;
        
        public:
            Parser(std::string& in_file_path, std::string& out_file_path);
            ~Parser();
        void StartWork();
    };
}

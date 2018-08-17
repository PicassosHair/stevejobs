#include <iostream>
#include <fstream>
#include <ctime>
#include "parser.hpp"

int main(int argc, char** argv) {
    
    clock_t begin = clock();
    
    std::ifstream inFile;
    std::ofstream outFile ("/Users/hao/Desktop/data/out.txt");
    inFile.open("/Users/hao/Desktop/data/1964.json");
    
    
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
    
    clock_t end = clock();
    double elapsed_secs = double(end - begin) / CLOCKS_PER_SEC;
    
    std::cout << "Seconds used: " << elapsed_secs;
    
    return 0;
}

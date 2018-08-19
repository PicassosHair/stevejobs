#include <iostream>
#include <ctime>
#include "parser.hpp"

using namespace idsguard;

int main(int argc, char** argv) {
    
    clock_t begin = clock();
    std::string in_file = "/Users/hao/Desktop/data/1964.json";
    Parser parser(in_file);
    
//    parser.StartWork();
    parser.ConnectDB();
    
    clock_t end = clock();
    double elapsed_secs = double(end - begin) / CLOCKS_PER_SEC;
    
    std::cout << "Seconds used: " << elapsed_secs;
    
    return 0;
}

#include <iostream>
#include <ctime>
#include "parser.hpp"

using namespace idsguard;

int main(int argc, char** argv) {
    
    clock_t begin = clock();
    std::string in_file = "/Users/hao/Desktop/data/1964.json";
    std::string out_file = "/Users/hao/Desktop/data/out.txt";
    Parser parser(in_file, out_file);
    
    parser.StartWork();
    
    clock_t end = clock();
    double elapsed_secs = double(end - begin) / CLOCKS_PER_SEC;
    
    std::cout << "Seconds used: " << elapsed_secs;
    
    return 0;
}

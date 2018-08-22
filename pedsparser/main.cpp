#include <iostream>
#include <ctime>
#include "parser.hpp"
#include "database.hpp"

using namespace idsguard;
using namespace std;

int main(int argc, char** argv) {
    const clock_t begin_time = clock();
    const string in_file = "/Users/hao/Desktop/data/1964.json";
    const string host_name = "tcp://127.0.0.1:3306/";
    const string username = "root";
    const string password = "";
    const string database_name = "idsguard_dev";
    
    Database database(host_name, username, password, database_name);
    Parser parser(in_file, &database);
    
    int work_status = parser.StartWork();
    if (work_status != 0) {
        std::cout << "Parser cannot work." << '\n';
        return -1;
    }
    
    clock_t end_time = clock();
    double elapsed_secs = double(end_time - begin_time) / CLOCKS_PER_SEC;
    
    std::cout << "Seconds used: " << elapsed_secs << '\n';
    
    return 0;
}

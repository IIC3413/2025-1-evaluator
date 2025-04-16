#include "util.h"
#include <string>

// Idea:
// Insert records in excess of the capacity of a single page to determine if submission is both able to
// properly add records to a page but also determine when there is no space available in one.

void run(std::string table_name) {
  for (int i = 0; i < 1000; i++) {
    catalog.insert_record(table_name, {"class_" + std::to_string(i), i, "name_" + std::to_string(i)});
  }
}

int main() {
  auto test_name = std::string("insert_1");
  auto system = System::init(Utils::get_db_name(test_name), BUFF_SIZE);
  Utils::run_test(test_name, run);
}

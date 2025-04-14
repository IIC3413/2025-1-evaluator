#include "util.h"
#include <string>

// Idea:
// Create a page with deleted records that has all its space used in order to see if the submission is able to
// recognize that even if an index is free it may not be usable.

void run(std::string table_name) {
  catalog.insert_record(table_name, {"s_class_" + std::to_string(0), 0, "s_name_" + std::to_string(0)});
}

int main() {
  auto test_name = std::string("insert_3");
  auto system = System::init(Utils::get_db_name(test_name), BUFF_SIZE);
  Utils::run_test(test_name, run);
}

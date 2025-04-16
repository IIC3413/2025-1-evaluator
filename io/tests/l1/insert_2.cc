#include "util.h"
#include <string>

// Idea:
// Create a page that has three deleted records (on both edges and in the center) in order to see if the
// submission is able to recognize that an index is reusable and use it appropriately.

void run(std::string table_name) {
  for (int i = 0; i < 3; i++) {
    catalog.insert_record(table_name, {"s_class_" + std::to_string(i), i, "s_name_" + std::to_string(i)});
  }
}

int main() {
  auto test_name = std::string("insert_2");
  auto system = System::init(Utils::get_db_name(test_name), BUFF_SIZE);
  Utils::run_test(test_name, run);
}

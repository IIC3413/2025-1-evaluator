#include "util.h"
#include <string>

// Idea:
// Evaluate if the submission is able to properly update the page's information and maintain the remaining
// record's order when vacuuming.

void run(std::string table_name) {
  catalog.get_table_info(table_name).heap_file->vacuum();
  Utils::zero_out_pages_free_space(table_name);
}

int main() {
  auto test_name = std::string("vacuum_1");
  auto system = System::init(Utils::get_db_name(test_name), BUFF_SIZE);
  Utils::run_test(test_name, run);
}

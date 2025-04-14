#include "util.h"
#include <string>

// Idea:
// Evaluate if the submission is able to delete every other record.

void run(std::string table_name) {
  Record buf(*catalog.get_table_info(table_name).schema);
  auto rids = std::vector<RID>();
  auto table_iter = catalog.get_table_info(table_name).heap_file->get_record_iter();

  table_iter->begin(buf);
  while (table_iter->next()) {
    if (buf.values[1].value.as_int % 2 != 0) {
      continue;
    }
    rids.push_back(table_iter->get_current_RID());
  }

  for (const RID& rid : rids) {
    catalog.delete_record(table_name, rid);
  }
}

int main() {
  auto test_name = std::string("delete_1");
  auto system = System::init(Utils::get_db_name(test_name), BUFF_SIZE);
  Utils::run_test(test_name, run);
}

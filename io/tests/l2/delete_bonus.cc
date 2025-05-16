#include "util.h"
#include <string>

// Idea:
// This test evaluates if a submission is abler to split a directory at the middle of a directory's pointers.

void run(std::string test_name) {
  Schema schema(COLUMNS);
  auto iter = catalog.get_table(test_name, &schema)->get_record_iter();
  Record buf(schema);

  iter->begin(buf);
  auto to_delete = std::vector<RID>();
  while (iter->next()) {
    if (buf.values[1].value.as_int % 2 != 0) {
      to_delete.push_back(iter->get_current_RID());
    }
  }
  for (auto rid : to_delete) {
    catalog.delete_record(test_name, rid);
  }
}

int main() {
  auto test_name = std::string("delete_bonus");
  // auto system = Utils::init_system(test_name);
  Utils::run(test_name, run, Utils::clone_output);
}

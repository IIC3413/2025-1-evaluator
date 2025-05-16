#include "storage/b_plus_tree/b_plus_tree_leaf.h"
#include "util.h"
#include <string>

// Idea:
// This test evaluates if a submission is able to insert a record at the end of a leaf.

void run(std::string test_name) {
  catalog.insert_record(
      test_name, {"class_" + std::to_string(BPlusTreeLeaf::max_records), BPlusTreeLeaf::max_records,
                  "name_" + std::to_string(BPlusTreeLeaf::max_records)}
  );
}

int main() {
  auto test_name = std::string("simple_insert_2");
  // auto system = Utils::init_system(test_name);
  Utils::run(test_name, run, Utils::clone_output);
}

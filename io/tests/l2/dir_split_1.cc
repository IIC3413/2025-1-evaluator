#include "storage/b_plus_tree/b_plus_tree_dir.h"
#include "storage/b_plus_tree/b_plus_tree_leaf.h"
#include "util.h"
#include <string>

// Idea:
// This test evaluates if a submission is abler to split a directory at the end of a directory pointers.

void run(std::string test_name) {
  auto base = BPlusTreeDir::max_children * BPlusTreeLeaf::max_records / 2 - 1;
  for (int i = 0; i < BPlusTreeLeaf::max_records + 1; i++) {
    auto is = std::to_string(base + i);
    catalog.insert_record(test_name, {"class_" + is, i, "name_" + is});
  }
}

int main() {
  auto test_name = std::string("dir_split_1");
  // auto system = Utils::init_system(test_name);
  Utils::run(test_name, run, Utils::clone_output);
}

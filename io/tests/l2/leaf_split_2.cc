#include "storage/b_plus_tree/b_plus_tree_leaf.h"
#include "util.h"
#include <string>

// Idea:
// This test evaluates if a submission is able to split a leaf not at the end of a directory pointers.

void run(std::string test_name) {
  for (int i = 0; i < BPlusTreeLeaf::max_records + 2; i++) {
    if (i % 2 == 0) {
      continue;
    }
    auto is = std::to_string(i);
    catalog.insert_record(test_name, {"class_" + is, i, "name_" + is});
  }
}

int main() {
  auto test_name = std::string("leaf_split_2");
  // auto system = Utils::init_system(test_name);
  Utils::run(test_name, run, Utils::clone_output);
}

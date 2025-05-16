#include "util.h"
#include <string>

// Idea:
// This test evaluates if a submission is able to split a leaf at the end of a directory, that is, without
// needing to move any of its pointers.

void run(std::string test_name) {
  catalog.insert_record(test_name, {"class_" + std::to_string(-1), -1, "name_" + std::to_string(-1)});
}

int main() {
  auto test_name = std::string("leaf_split_1");
  // auto system = Utils::init_system(test_name);
  Utils::run(test_name, run, Utils::clone_output);
}

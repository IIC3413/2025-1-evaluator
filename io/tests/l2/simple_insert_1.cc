#include "util.h"
#include <string>

// Idea:
// We are just being nice here. Simple insert on an empty index.

void run(std::string test_name) {
  catalog.insert_record(test_name, {"class_5", 5, "name_5"});
}

int main() {
  auto test_name = std::string("simple_insert_1");
  Utils::run(test_name, run, Utils::clone_output);
}

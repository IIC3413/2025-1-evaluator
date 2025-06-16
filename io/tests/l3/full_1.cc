#include "util.h"
#include <string>

// Idea:
// Simple full outer join. We aim to evaluate if submissions are able to
// match tuples and correctly handle transitioning from a no match state to
// a match one.

void setup() {
  Utils::create_table("ltable", {{"la", DataType::STR}, {"lb", DataType::INT}});
  Utils::create_table("rtable", {{"ra", DataType::STR}, {"rb", DataType::INT}});
  for (int i = 0; i < 10; i++) {
    catalog.insert_record("ltable", {std::to_string(i), i});
  }
  for (int i = 5; i < 15; i++) {
    catalog.insert_record("rtable", {std::to_string(i), i});
  }
}

int main() {
  auto test_name = std::string("full_1");
  auto query = std::string("SELECT * FROM (ltable FULL OUTER JOIN rtable ON ltable.lb == rtable.rb) AS FOJ1");
  Utils::run(test_name, setup, query);
}

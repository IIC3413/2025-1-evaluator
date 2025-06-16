#include "relational_model/schema.h"
#include "util.h"
#include <string>

// Idea:
// Same as full_1 but we incorporate a second criteria.

void setup() {
  Utils::create_table("ltable", {{"la", DataType::STR}, {"lb", DataType::INT}});
  Utils::create_table("rtable", {{"ra", DataType::STR}, {"rb", DataType::INT}});
  for (int i = 0; i < 15; i++) {
    catalog.insert_record("ltable", {i % 2 == 0 ? "0" : "1", i});
  }
  for (int i = 5; i < 20; i++) {
    catalog.insert_record("rtable", {i % 2 == 0 ? "0" : "1", i});
  }
}

int main() {
  auto test_name = std::string("full_2");
  auto query = std::string(
      "SELECT * FROM (ltable FULL OUTER JOIN rtable ON ltable.lb == rtable.rb AND ltable.la == rtable.ra) AS FOJ2"
  );
  Utils::run(test_name, setup, query);
}


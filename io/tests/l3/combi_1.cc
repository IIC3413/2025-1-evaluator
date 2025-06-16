#include "util.h"
#include <string>

// Idea:
// Just looking for correct implementation of the query iterator interface.

void setup() {
  Utils::create_table("ltable1", {{"l1a", DataType::INT}, {"l1b", DataType::INT}});
  Utils::create_table("rtable1", {{"r1a", DataType::INT}, {"r1b", DataType::INT}});
  Utils::create_table("ltable2", {{"l2a", DataType::STR}, {"l2b", DataType::INT}});
  Utils::create_table("rtable2", {{"r2a", DataType::STR}, {"r2b", DataType::INT}});
  for (int i = 0; i < 10; i++) {
    catalog.insert_record("ltable1", {i * 4, i * 2});
  }
  for (int i = 5; i < 15; i++) {
    catalog.insert_record("rtable1", {i * 2, i});
  }
  for (int i = 0; i < 15; i++) {
    catalog.insert_record("ltable2", {std::to_string(i), i});
  }
  for (int i = 5; i < 20; i++) {
    catalog.insert_record("rtable2", {std::to_string(i), i});
  }
}

int main() {
  auto test_name = std::string("combi_1");
  auto query = std::string("SELECT * FROM \
    (ltable1 LEFT OUTER JOIN rtable1 ON ltable1.l1b == rtable1.r1b AND ltable1.l1a == rtable1.r1a) AS LOJ1, \
    (ltable2 FULL OUTER JOIN rtable2 ON ltable2.l2a == rtable2.r2a) AS FOJ1 \
    WHERE LOJ1.r1b == FOJ1.l2b");
  Utils::run(test_name, setup, query);
}

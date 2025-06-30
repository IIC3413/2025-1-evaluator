#include "util.h"
#include <string>

void setup() {
  Utils::create_table("A", {{"x", DataType::STR}, {"y", DataType::INT}, {"z", DataType::INT}});
  Utils::create_table("B", {{"x", DataType::STR}, {"y", DataType::INT}});
  for (int i = 0; i < 100; i++) {
    catalog.insert_record("A", {std::to_string(i), i, i * 2});
  }
  for (int i = 50; i < 150; i++) {
    catalog.insert_record("B", {std::to_string(i), i});
  }
}

int main() {
  auto test_name = std::string("p1");
  auto query = std::string("\
    SELECT * \
    FROM A, B \
    WHERE A.x == B.x AND A.y == B.y");
  Utils::run_bloom_filter(test_name, setup, query);
}

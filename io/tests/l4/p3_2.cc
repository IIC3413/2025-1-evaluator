#include "util.h"
#include <string>

void setup() {
  Utils::create_table("A", {{"x", DataType::STR}, {"y", DataType::INT}, {"z", DataType::INT}, {"w", DataType::STR}});
  Utils::create_table("B", {{"x", DataType::STR}, {"y", DataType::INT}, {"z", DataType::INT}});
  for (int i = 0; i < 200; i++) {
    catalog.insert_record("A", {std::to_string(i % 100), i % 100, i, std::to_string(i % 2)});
  }
  for (int i = 50; i < 150; i++) {
    catalog.insert_record("B", {std::to_string(i), i, i % 2});
  }
}

int main() {
  auto test_name = std::string("p3_2");
  auto query = std::string("\
    SELECT * \
    FROM A, B \
    WHERE A.x == B.x AND A.z == B.y");
  Utils::run(test_name, setup, query);
}

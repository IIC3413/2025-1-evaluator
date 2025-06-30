#include "util.h"
#include <string>

void setup() {
  Utils::create_table("A", {{"x", DataType::INT}, {"y", DataType::INT}});
  Utils::create_table("B", {{"x", DataType::INT}, {"y", DataType::INT}});
  Utils::create_table("C", {{"x", DataType::INT}, {"y", DataType::INT}});
  Utils::create_table("D", {{"x", DataType::INT}, {"y", DataType::INT}});
  for (int i = 0; i < 50; i++) {
    catalog.insert_record("A", {i, i + 1});
  }
  for (int i = 0; i < 200; i++) {
    catalog.insert_record("B", {i * 2, i});
  }
  for (int i = 0; i < 150; i++) {
    catalog.insert_record("C", {i, i % 25});
  }
  for (int i = 0; i < 1000; i++) {
    catalog.insert_record("D", {i, i % 75});
  }
}

int main() {
  auto test_name = std::string("p4_2");
  auto query = std::string("\
    SELECT * \
    FROM A, B, C, D \
    WHERE A.y == B.x \
    AND B.y == C.y \
    AND C.x == D.y");
  Utils::run(test_name, setup, query);
}

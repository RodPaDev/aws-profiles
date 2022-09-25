import test from "ava";
import { createSpaces } from "../../src/utils";

test("create 2 spaces", (t) => {
  t.is("  ", createSpaces(2));
});

test("create 1 space", (t) => {
  t.is(" ", createSpaces(1));
});

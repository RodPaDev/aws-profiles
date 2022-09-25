import { createSpaces } from "../../src/util";

test("create one space", () => {
  expect(createSpaces()).toBe(" ");
});

test("create one space", () => {
  expect(createSpaces()).toBe(" ");
});

test("create 2 spaces", () => {
  expect(createSpaces(2)).toBe("  ");
});

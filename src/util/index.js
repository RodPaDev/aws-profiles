export function createSpaces(spaces = 1) {
  let str = "";
  let idx = 0;
  while (idx < spaces) {
    str += " ";
    idx += 1;
  }
  return str;
}

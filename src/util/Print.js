import chalk from "chalk";

export class Print {
  /* eslint-disable no-console */
  static log(...args) {
    console.log(...args);
  }
  static error(...args) {
    console.error(chalk.redBright("Error: ", ...args));
  }
  static info(...args) {
    console.info(...args);
  }
  static warn(...args) {
    console.warn(...args);
  }
  /* eslint-enable no-console */
}

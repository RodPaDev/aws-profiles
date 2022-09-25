import chalk from "chalk";

export class Print {
  /* eslint-disable no-console */
  static log(...args: any[]) {
    console.log(...args);
  }
  static error(...args: any[]) {
    console.error(chalk.redBright("Error: ", ...args));
  }
  static info(...args: any[]) {
    console.info(...args);
  }
  static warn(...args: any[]) {
    console.warn(...args);
  }
  /* eslint-enable no-console */
}

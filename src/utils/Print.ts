import chalk from 'chalk'

export class Print {
  /* eslint-disable no-console */
  static log(...args: unknown[]) {
    console.log(...args)
  }
  static error(...args: unknown[]) {
    console.error(chalk.redBright('Error: ', ...args))
  }
  static info(...args: unknown[]) {
    console.info(...args)
  }
  static warn(...args: unknown[]) {
    console.warn(...args)
  }
  /* eslint-enable no-console */
}

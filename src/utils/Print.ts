import chalk from 'chalk'

export class Print {
  /* eslint-disable no-console */
  static log(...args: unknown[]) {
    console.log(...args)
  }
  static error(...args: unknown[]) {
    console.error(chalk.redBright('x'), chalk.white(...args))
  }
  static info(...args: unknown[]) {
    console.info(chalk.green('\u2713'), chalk.white(...args))
  }
  static warn(...args: unknown[]) {
    console.warn(chalk.yellow('!'), chalk.white(...args))
  }
  /* eslint-enable no-console */
}

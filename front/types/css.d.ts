// This file overrides type of *.module.css modules for tests
// Since ts-jest for some reason ignores typescript's compiler plugins

declare module "*.module.css" {
  const CSSModule: {
    [key: string]: string;
  }

  export default CSSModule
}

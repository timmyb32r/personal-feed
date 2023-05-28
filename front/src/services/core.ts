export function wrapURI(uri: string): string {
    if (process.env.NODE_ENV == "development") {
        return "http://localhost".concat(uri)
    }
    return uri
}

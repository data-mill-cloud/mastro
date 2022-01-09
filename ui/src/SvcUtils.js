export const getSvcHost = (svcName) => {
    const devEnvVar = `REACT_APP_${svcName.toUpperCase()}_URL`
    return typeof process.env[devEnvVar] !== 'undefined' ? process.env[devEnvVar] : svcName
}
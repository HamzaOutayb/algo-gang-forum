


async function isLoggedIn() {
    const response = await fetch("/api/checkuser").catch((e) => { console.log(e) })
    if (!response.ok) {
        return
    }
    return await response.json()
}

export default { isLoggedIn }
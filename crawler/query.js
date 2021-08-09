async function query(page, dateRange) {
  await page.goto("https://github.com/trending/go?since=" + dateRange, { waitUntil: "load" })

  let nameAndUrl = await page.$$eval(".Box-row > h1 > a", item => item.map(e => { return { name: e.text, href: e.href } }))
  let overview = await page.$$eval(".Box-row > p", item => item.map(e => e.innerHTML))
  let star = await page.$$eval(".Box-row > .f6 > span + a", item => item.map(e => e.text))
  let fork = await page.$$eval(".Box-row > .f6 > span + a + a", item => item.map(e => e.text))
  let todayStar = await page.$$eval(".Box-row > .f6 > .float-sm-right", item => item.map(e => e.innerHTML))

  let temp = []
  let boxCount = await page.$$eval(".Box-row", item => item.length)
  for (let i = 0; i < boxCount; i++) {
    temp.push({
      "name": nameAndUrl[i].name.replace(/[ ]|[\r\n]/g, ""),
      "url": nameAndUrl[i].href.replace(/[ ]|[\r\n]/g, ""),
      "overview": overview[i].replace(/<[^>]+>|[\r\n]/g, "").trim(),
      "star": +star[i].replace(/[^\d.]/g, ""),
      "todayStar": +todayStar[i].replace(/<[^>]+>|[\r\n]|[^\d.]/g, ""),
      "fork": +fork[i].replace(/[^\d.]/g, ""),
    })
  }

  return {
    date: new Date(),
    [dateRange]: temp.sort((a, b) => b.star - a.star)
  }
}

module.exports = { query }
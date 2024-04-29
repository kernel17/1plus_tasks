document.addEventListener("DOMContentLoaded", main);

async function main() {
    try {
        const response = await fetch("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1");
        if (!response.ok) {
          throw new Error("error while getting currencies");
        }
        const json = await response.json();
        const table = document.getElementById("coins")
        let counter = 0
        json.map(e => {
            const row = table.insertRow(-1)
            const c1 = row.insertCell(0)
            const c2 = row.insertCell(1)
            const c3 = row.insertCell(2)
            c1.innerText = e.id
            c2.innerText = e.symbol
            c3.innerText = e.name
            if (counter < 5) {
                c1.classList.add("blue")
                c2.classList.add("blue")
                c3.classList.add("blue")
                counter++
            }
            if (e.symbol == "usdt") {
                c1.classList.add("bucks")
                c2.classList.add("bucks")
                c3.classList.add("bucks")
            }
        })
      } catch (error) {
        console.log(error.message);
      }
   
}
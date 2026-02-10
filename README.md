# MimeaLogic 🌿

**SDG 6: Clean Water and Sanitation | Sustainable Agricultural Intelligence**

MimeaLogic is a high-efficiency irrigation controller built in Golang. It utilizes an **O(1) Math Engine** to predict soil moisture decay without the need for high-energy constant sensor polling.

## 🚀 Key Features
- **Stateless Prediction:** Uses logarithmic decay math to estimate soil moisture.
- **Persistence:** Maintains state using a flat-file ledger (`last_watered.txt`).
- **Energy Optimized:** Low-duty cycle agent designed for off-grid solar hardware.
- **Strictly Typed:** Built with Go for maximum reliability and memory safety.

## 🛠 Project Structure
- `cmd/main.go`: The central control agent and infinite loop.
- `pkg/engine.go`: The core mathematical logic (The "Brain").
- `last_watered.txt`: The system's memory (Automated).

## 💻 Installation & Usage
1. Clone the repository:
   `git clone https://github.com/xbuyan/mimealogic.git`
2. Run the agent:
   `go run cmd/main.go`

## 🌍 UN SDG Impact
By optimizing water usage through predictive math rather than wasteful constant irrigation, MimeaLogic directly supports **UN Sustainable Development Goal 6.4**: *Substantially increase water-use efficiency across all sectors.*
Imagine you're working on a web application on your personal computer. Normally, this app is only accessible on your local network (localhost). A tunneling tool is like a magical bridge that lets you
- Expose your local server to the entire internet
- Create a public URL for your local development
- Share your work without complex network configurations

The Problem Tunneling Solves:
Let's say you're building a cool web app:

- You run it on localhost:3000
- Only you can see it on your computer
- Your friend or potential client can't access it

Tunneling tools solve this by:

- Creating a public endpoint
- Forwarding all internet traffic to your local server
- Allowing anyone to access your app via a temporary public URL

But how do you actually access it on external devices like your phone/tablets

You need to use the actual IP address of the machine running the tunnel.
so the forwarding Url it spits out is basically http://<your local machines IP><random public port>

Firewall/Network Configuration:

Ensure your machine's firewall is not blocking incoming connections on the generated port.
Check that your network router is not preventing external access.

using this is very simple

clone the repo
- run go run <filename.go>

- start any local server you have on port 3000 (if you use a different port number you can run go run <filename.go> PORT)
- OUTPUT: 
🌐 Tunnel created successfully!
📍 Local Server: localhost:3004
🌍 Public Endpoint: 192.-.-.-:54321

app-aurora
==========

A Ninja Sphere app that can help you photograph an aurora!

Introduction
============

Aurora hunters use four key metrics for determining aurora activity and potential visibility -- Kp (or Planetary Index), Density, Speed, and Bz. When those metrics reach certain levels, grab your camera and head towards your nearest polar region, because you could see the northern (or southern) lights.

I, like many others, often rely on websites such as http://aurora-service.net to get the information we need. However logging on to the site and scrolling for data is so 1996, so I wrote a Ninja Sphere app that builds upon [go-aurora][1] to grab information directly from NOAA's Space Weather Prediction Center (which grabs its data from the ACE spacecraft) and display it in a user-friendly format. Now if I want to know if I should get out of bed and grab my camera, I can just roll over and swipe to the aurora pane for an instant look.

Installation & Usage
====================

A proper installer will be coming soon, but until then, you'll need to build the binary yourself (`GOOS=linux GOARCH=arm go build`) and copy it, and `package.json` into a folder called `app-aurora` in the  `/data/sphere/user-autostart/apps` folder

Usage is simple. Swipe to the new pane to see the gauge. Green = no activity, red = lots of activity. Tap the Sphere multiple times to cycle through the four points of dataTimeout

To-Do
=====

- [ ] Change the Kp gauge to a score gauge. Kp can be a poor indicator of aurora activity
- [ ] Colour code the values shown when you tap
- [ ] Overlay the Kp on the gauge
- [ ] Perhaps investigate scrolling to allow larger data (e.g. speed) to be accurately displayed.

Helping out
===========

The base library (go-aurora) uses a weighted score. For example, if Speed is in the red (that is, higher than 700km/s), then 25 points are added to the score, while a green result (e.g. Between 200-350 km/s) will get -10 points, and no data gets even more points taken off. This system isn't perfect, and I need your help to make it better. Speed, Bz and Density are "worth more" than Kp, given that I witnessed a Kp 5 storm that was bright, and a Kp 7 storm that was barely visible, so pull requests and such are most welcome.  


[1]: http://github.com/Grayda/go-aurora

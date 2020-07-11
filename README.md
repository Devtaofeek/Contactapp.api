# Contactapp.api
This is a secured restful web api written in Golang, It saves,updates,creates and deletes contacts from a postgres database.
The intial code base for this project was written by https://github.com/adigunhammedolalekan, I learnt how to build web api using go programming language by studying this project,

I rewrote the code three times until i understood what every line does.
This repo is the last rewrite and alot of things have changed compared to the code base i learnt from.

I made sure  integrity was enforced in the database by implenting the appropriate database relatioships,

I added refresh tokens in other for better security and improved user experience, refresh tokens enables issuance of another access token without logging the user out,

I added the use of reddis for persistence storage of the access and refresh tokens, reddis has a functionality that automatically deletes the access and refresh tokens upon expiry,

reddis also takes care of a scenario where a user logs out while the access token is still valid, to prevent any security breach reddis is used to invalidate the token when a user logs out ,

I made use of uuid as the keys in reddis store, this uuid is also part of the access and refresh token claims.
I added  some additional end points for managing contacts.
I hope to still improve more on this project as i continue to learn Go.

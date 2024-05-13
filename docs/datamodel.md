# Data Model

I don't think this bot can be completely stateless - we might need to keep the session alive for button interactions to be persistent.

We want to store the following things somewhere:

- what challenges have been generated, so we have a reference for the following
- what completion times are marked against that challenge, so we have a record of how well we did against each other
- whether the car / track / weather combination for a particular challenge was good or bad, to modify future challenge generation
- the current state of the challenge generation modifiers

That translates to the following set of events:

- Creation
  - challenge details
- Completion
  - user ID
  - user display name
  - challenge ID
  - duration
- Feedback
  - user ID
  - challenge ID
  - challenge feedback

and some sort of lump of modifiers - for now, these will live in memory, but in the interest of making everything stateless they could go and live with whatever persistent store the events use.

Discord uses unique IDs as documented here: https://discord.com/developers/docs/reference#snowflakes

This provides a handy event ID - we can use the ID of interaction that causes an event to be created.

The snowflake also contains a timestamp to the millisecond level, which is good enough for our purposes - when (re)hydrating the bot we can use the snowflake to provide an event order. While imperfect, the only case where order really matters is ensuring events relating to a particular challenge being processed after the challenge creation event.

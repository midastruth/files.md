
# Bot's State Machine

It's just a subset of the whole state machine

```plantuml
@startuml
    'skinparam linetype polyline
    'skinparam linetype ortho
    
    state cmd <<start>>
    state cmdTypes <<choice>>
    state "Added to Default Dir" as added
    state "Day Selector" as daySelectorOnce
    state "Today" as today
    state "Later" as later

    
    
    [*] --> added: User typed arbitrary text
    cmd --> cmdTypes : User entered command
    cmdTypes --> later : /later
    cmdTypes --> today : /today
    added --> daySelectorOnce : "For a Day" clicked
    added --> today : "Today" clicked
    added --> today : "For tmrw" clicked
    added --> later : "For later" clicked
    daySelectorOnce --> daySelectorRecurrent : "Repeat the task" clicked
    daySelectorOnce --> today : "Cancel" clicked
    daySelectorRecurrent -> today : Some day selected
    
    later --> today : "Tasks for today" clicked
    today --> later : "Tasks for later" clicked
    today --> today : Oneline task clicked
    later --> later : Oneline task clicked
    today  --> multiline_task : Multiline task clicked
    later --> multiline_task : Multiline task clicked
    
    state "Multiline Task Processing" as multiline_task {
        state "Multiline Task Shown" as ml
        state "Get back to today/later" as back
        [*] -> ml
        ml --> back : "Back" clicked
        ml --> back : "Complete" clicked
        ml --> back : "To Later/Today" clicked
        back --> [*]    
    }

@enduml
```

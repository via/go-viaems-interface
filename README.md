== Usage
This provides an HTTP interface to a viaems target.  The interface allows access
to configuration state, as well as current target status and history.

The target can be a real target (hardware or simulation) over uart or tcp, or it
can be a virtual target from logs for replay.  In the latter, a virtual
configuration model can be loaded from file for the log.

If the target is real and is getting status updates, they will be recorded to
the selected files.

== Backend API

/info

/target/config

/target/status
/target/status?start_time=blah&stop_time=blah

GET -> {
  "sensors": {
    "MAP": {
      "value": 12.2,
      "fault": null
    }
  },
  "fueling": {
    "ve": 50.0,
    "pulsewidth": 1200,
  },
  "ignition": {
    "advance": 12.0,
  },
  "decoder": {
    "t0": 1
  },
  "cpu_time": 1234,
  "real_time": 6789
}

POST {
  "sensors": {
    "EGO": {
      "value": true,
    }
  }
-> 
  "sensors": {
    "EGO": {
      "value": 12.2
    }
  }


/target/meta
GET -> {
  "type": "uart",
  "history_state": {
    "file": "/home/via/blah",
    "duration": 92.2
  },
  "link_state": "connected | disconnected",
  "model_state": "synchronized"
}

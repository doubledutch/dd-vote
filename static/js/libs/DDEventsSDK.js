// Copyright © DoubleDutch 2015

(function () {
  function throwNotInitialized() {
    throw new Error("DD Events SDK runtime not initialized");
  }

  // This is a simple queue to handle our calls into native code
  // on the mobile apps. The problem is that the calls are async
  // and we were using a single callback variable. If you made multiple
  // calls to one of the APIs below before each had responded,
  // the last callback passed would be the only one called, but
  // with the result of all the previous calls
  var CallbackQueue = function() {
    var queue = []
    var activeCallback = throwNotInitialized

    this.callback = function() {
      activeCallback.apply(undefined, arguments)
      if (queue.length) {
        var nextCall = queue.splice(0, 1)[0]
        queue = nextCall[0]
        nextCall[1].apply(undefined, arguments)
      }
    }

    this.enqueue = function (callback, operation) {
      if (queue.length) {
        // Enqueue the request
        queue.push([callback, operation])
      } else {
        // Nothing is queued - Go ahead and call it
        activeCallback = callback
        operation()
      }
    }
  }

  var onReadyCallbacks = [];
  var DD = {
    Events: {
      // Subscribe by passing a callback to onReady.
      // Callbacks will be called after the DD Events Platform is initialized.
      onReady: function (callback) {
        if (onReadyCallbacks) {
          onReadyCallbacks.push(callback);
        } else {
          setTimeout(callback, 1);
        }
      },

      /// Get Current User
      /// The implementation is set by the hosting native application
      getCurrentUserImplementation: throwNotInitialized,
      _getCurrentUserQueue: new CallbackQueue(),
      /// This async method is called by JS clients
      getCurrentUserAsync: function (callback) {
        var self = this
        this.getCurrentUserCallback = this._getCurrentUserQueue.callback
        this._getCurrentUserQueue.enqueue(callback, function() { self.getCurrentUserImplementation() })
      },

      /// Get Current Event
      /// The implementation is set by the hosting native application
      getCurrentEventImplementation: throwNotInitialized,
      _getCurrentEventQueue: new CallbackQueue(),
      /// This async method is called by JS clients
      getCurrentEventAsync: function (callback) {
        var self = this
        this.getCurrentEventCallback = this._getCurrentEventQueue.callback
        this._getCurrentEventQueue.enqueue(callback, function() { self.getCurrentEventImplementation() })
      },

      /// Get OAuth Encoded API Call Implementation
      /// The implementation is set by the hosting native application
      getSignedAPIImplementation: throwNotInitialized,
      _getSignedAPIQueue: new CallbackQueue(),
      /// This async method is called by JS clients
      getSignedAPIAsync: function (apiFragment, postBody, callback) {
        var self = this
        this.getSignedAPICallback = this._getSignedAPIQueue.callback
        this._getSignedAPIQueue.enqueue(callback, function() { self.getSignedAPIImplementation(apiFragment, postBody) })
      },

      /// Set an action button in the native app (if implemented)
      /// The implementation is set by the hosting native application
      setActionButtonImplementation: throwNotInitialized,
      _getActionButtonQueue: new CallbackQueue(),
      /// This async method is called by JS clients
      setActionButtonAsync: function (title, imageReserved, callback) {
        var self = this
        this.getActionButtonCallback = this._getActionButtonQueue.callback
        this._getActionButtonQueue.enqueue(callback, function() { self.getActionButtonImplementation(title, imageReserved) })
      },

      /// Set the title in the native app (if implemented)
      /// The implementation is set by the hosting native application
      setTitleImplementation: throwNotInitialized,
      setTitleAsync: function (title) {
        this.setTitleImplementation(title);
      }
    }
  };

  var initCheck = setInterval(function () {
    if (DD.Events.getSignedAPIImplementation !== throwNotInitialized) {
      clearInterval(initCheck);
      var cbs = onReadyCallbacks;
      onReadyCallbacks = null;
      for (var i = 0; i < cbs.length; ++i) {
        try{
          cbs[i]();
        } catch (e) { console.log(e); }
      }
    }
  }, 25);

  window.DD = DD;

  // OAuth2-compatible implementation for HTML5.
  window.addEventListener("message", function(event)
  {
    var data = event.data;

    DD.Events.setActionButtonImplementation = function() {};
    DD.Events.setTitleImplementation = function(title) {
      window.document.title = title;
    };

    if (data.user) {
      var currentUser = data.user;
      DD.Events.getCurrentUserImplementation = function() {
        DD.Events.getCurrentUserCallback(currentUser);
      };
    }

    if (data.event) {
      var currentEvent = data.event;
      DD.Events.getCurrentEventImplementation = function() {
        DD.Events.getCurrentEventCallback(currentEvent);
      };
    }

    if (data.authorizationHeader && data.apiRoot) {
      var auth = data.authorizationHeader;
      var apiRoot = data.apiRoot;

      DD.Events.getSignedAPIImplementation = function(apiFragment, postBody) {
        var url = apiRoot + apiFragment;
        url += (url.indexOf('?') < 0 ? '?' : '&') + 'sdk=true';
        DD.Events.getSignedAPICallback(url, auth);
      };
    }
  }, false);
})();

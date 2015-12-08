var controller = angular.module("controller", []);
controller.controller("controller", ["$scope", function ($scope) {

    $scope.tasks = [];
    $scope.updateHours;

    $scope.newTaskName = "";
    $scope.initialEstimate = "";

    $scope.init = function() {
        $scope.addTask("Sample Task 1", "", 8);
        $scope.addTask("Sample Task 2", "", 15);
        $scope.updateTask("Sample Task 2", 4);
    }

    $scope.addTask = function(name, description, hours) {
        if (name === "") {
            alert("Task must have a name");
            return;
        }
        if (hours === "" || isNaN(parseFloat(hours))) {
            alert("Task must have a valid initial time estimate");
            return;
        }
        for (var k = 0; k < $scope.tasks.length; k++) {
            if ($scope.tasks[k].name === name) {
                alert("Task name already exists");
                return;
            }
        }
        var newTask = {
            "name" : name,
            "description" : description,
            "percentDone" : 0,
            "initialGoal" : hours,
            "hoursSpent" : 0,
            "hoursRemaining" : hours,
            "hidden" : false
        }
        $scope.tasks.push(newTask);

        $scope.newTaskName = "";
        $scope.initialEstimate = "";
    }

    $scope.makeComplete = function(task) {
        task.percentDone = 100;
        task.hoursRemaining = 0;
        task.hidden = true;
        return task;
    }

    $scope.updateTask = function(name, hoursSpent) {
        if (hoursSpent < 0 || !hoursSpent || isNaN(parseFloat(hoursSpent))) {
            alert("Please enter a positive number of hours");
            return;
        }
        for (var k = 0; k < $scope.tasks.length; k++) {
            if ($scope.tasks[k].name === name) {
                var hr = Math.round(($scope.tasks[k].hoursRemaining - hoursSpent) * 100) / 100;
                $scope.tasks[k].hoursSpent = Math.round(($scope.tasks[k].hoursSpent + Number(hoursSpent)) * 100) / 100;
                if (hr <= 0) {
                    $scope.tasks[k].percentDone = 100;
                    $scope.tasks[k].hoursRemaining = 0;
                    return;
                } else {
                    $scope.tasks[k].hoursRemaining = hr;
                }
                pd = 100 * ($scope.tasks[k].hoursSpent / $scope.tasks[k].initialGoal);

                $scope.tasks[k].percentDone = Math.round(pd);
            }
        }
        $scope.updateHours = 0;
        $scope.remainingHours = 0;
    }

    $scope.deleteTask = function(name) {
        for (var k = 0; k < $scope.tasks[k].length; k++) {
            if ($scope.tasks[k].name === name) {
                $scope.tasks[k].splice(k, 1);
            }
        }
    }

}]);
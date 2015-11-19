var dogpark = new angular.module("dogpark", ["ngResource"]);

	dogpark.controller("MainCtl", ["$scope", "$resource",
		function($scope, $resource) {
			var Dog = $resource("/dogs/:id", {
				id: "@id"
			}, {});
	
			$scope.selected = null;
			$scope.dogs = Dog.query();
	
	
			$scope.get = function(dog) {
				Dog.get({
					id: dog.id
				}, function(data) {
					$scope.selected = data;
				})
			};
	
			$scope.remove = function(dog) {
				$scope.selected.$delete(function() {
					$scope.dogs = Dog.query();
				});
	
			};
	
			$scope.add = function() {
				var name = prompt("Enter dog name.")
				if (!name) {
					alert("Name can't be empty.")
					return;
				}
				var owner = prompt("Enter dog owner.")
				if (!owner) {
					alert("Owner can't be empty.")
					return;
				}
	
				var dog = new Dog();
				dog.name = name;
				dog.owner = owner;
				dog.$save(function() {
					$scope.dogs = Dog.query();
				})
			}
	
			$scope.update = function(dog) {
	
				var name = prompt("Enter dog name.", dog.name)
				if (!name) {
					alert("Name can't be empty.")
					return;
				}
				var owner = prompt("Enter dog owner.", dog.owner)
				if (!owner) {
					alert("Owner can't be empty.")
					return;
				}
				
				dog.name = name;
				dog.owner = owner;
				
				dog.$save(function() {
					$scope.dogs = Dog.query();
				})
	
			}
	
		}
	]);
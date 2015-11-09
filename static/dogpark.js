	var myapp = new angular.module("dogpark", ["ngResource"]);

	myapp.controller("MainCtl", ["$scope", "$resource", function($scope, $resource){

		var dog = $resource("/books/:id", {id: '@id'}, {});

		$scope.selected = null;

		$scope.list = function(idx){
			// Notice calls to Book are often given callbacks.
			Dog.query(function(data){
				$scope.Dogs = data;
				if(idx != undefined) {
					$scope.selected = $scope.dogs[idx];
					$scope.selected.idx = idx;
				}
			}, function(error){
				alert(error.data);
			});
		};

		$scope.list();

		$scope.get = function(idx){
			// Passing parameters to Book calls will become arguments if
			// we haven't defined it as part of the path (we did with id)
			Dog.get({id: $scope.dogs[idx].id}, function(data){
				$scope.selected = data;
				$scope.selected.idx = idx;
			});
		};

		$scope.add = function() {
			// I was lazy with the user input.
			var name = prompt("Enter the dog's name.");
			if(name == null){
				return;
			}
			var owner = prompt("Enter the dogs's owner.");
			if(owner == null){
				return;
			}
			// Creating a blank object means you can still $save
			var newDog = new Dog();
			newDog.name = name;
			newDog.owner = owner;
			newDog.$save();

			$scope.list();
		};

		$scope.update = function(idx) {
			var dog = $scope.dogs[idx];
			var name = prompt("Enter a new dog", dog.name);
			if(name == null) {
				return;
			}
			var owner = prompt("Enter a new owner", dog.author);
			if(owner == null) {
				return;
			}
			dog.title = title;
			dog.author = author;
			// Noticed I never created a new Book()?
			dog.$save();

			$scope.list(idx);
		};

		$scope.remove = function(idx){
			$scope.dogs[idx].$delete();
			$scope.selected = null;
			$scope.list();
		};
	}]);
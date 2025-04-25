package models

func Map[T any, R any](arr []T, f func(T) R) []R {
	appliedArr := make([]R, len(arr))
	for i, value := range arr {
		appliedArr[i] = f(value)
	}
	return appliedArr
}

func ToHorseDto(horse Horse) HorseDto {
	return HorseDto{
		ID:     horse.ID,
		Name:   horse.Name,
		Breed:  horse.Breed,
		Age:    horse.Age,
		Gender: horse.Gender,
	}
}

func ToHorseDtos(horses []Horse) []HorseDto {
	return Map(horses, ToHorseDto)
}

func ToUserDto(user User) UserDto {
	return UserDto{
		ID:       user.ID,
		Username: user.Username,
		Mail:     user.Mail,
		Age:      user.Age,
		Gender:   user.Gender,
		Country:  user.Country,
		Roles:    user.Roles,
		Horses:   ToHorseDtos(user.Horses),
	}
}

func UpdateUserToUser(uUser UpdateUser) User {
	return User{
		Username: uUser.Username,
		Mail:     uUser.Mail,
		Password: uUser.Password,
		Age:      uUser.Age,
		Gender:   uUser.Gender,
		Country:  uUser.Country,
	}
}

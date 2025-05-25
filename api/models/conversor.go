package models

import "api/api/utils"

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
	return utils.Map(horses, ToHorseDto)
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

func ToGroupDto(group Group) GroupDto {
	return GroupDto{
		ID:    group.ID,
		Name:  group.Name,
		Users: uint(len(group.Users)),
	}
}

func ToGroupDtos(groups []Group) []GroupDto {
	return utils.Map(groups, ToGroupDto)
}

func ToMessageDto(message Message) MessageDto {
	return MessageDto{
		Username: message.User.Username,
		Time:     message.CreatedAt.String(),
		Message:  message.Message,
	}
}

func ToMessageDtos(messages []Message) []MessageDto {
	return utils.Map(messages, ToMessageDto)
}

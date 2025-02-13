package types

type Post struct {
	Id           string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId       string   `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserInfo     UserInfo `json:"userInfo"`
	Caption      string   `protobuf:"bytes,3,opt,name=caption,proto3" json:"caption,omitempty"`
	CreatedAt    int64    `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt    int64    `protobuf:"varint,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Mentions     []string `protobuf:"bytes,6,rep,name=mentions,proto3" json:"mentions,omitempty"`
	Tags         []string `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	Media        []*Media `protobuf:"bytes,8,rep,name=media,proto3" json:"media,omitempty"`
	Location     string   `protobuf:"bytes,9,opt,name=location,proto3" json:"location,omitempty"`
	LikeCount    int32    `protobuf:"varint,10,opt,name=like_count,json=likeCount,proto3" json:"like_count,omitempty"`
	CommentCount int64    `protobuf:"varint,11,opt,name=comment_count,json=commentCount,proto3" json:"comment_count,omitempty"`
}

type Media struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FileUrl   string `protobuf:"bytes,2,opt,name=file_url,json=fileUrl,proto3" json:"file_url,omitempty"`
	PublicId  string `protobuf:"bytes,3,opt,name=public_id,json=publicId,proto3" json:"public_id,omitempty"`
	FileType  string `protobuf:"bytes,4,opt,name=file_type,json=fileType,proto3" json:"file_type,omitempty"`
	FileName  string `protobuf:"bytes,5,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	CreatedAt int64  `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	PostId    string `protobuf:"bytes,7,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
}

type CreatePostRequest struct {
	UserId   string         `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Caption  string         `protobuf:"bytes,2,opt,name=caption,proto3" json:"caption,omitempty"`
	Location string         `protobuf:"bytes,3,opt,name=location,proto3" json:"location,omitempty"`
	Tags     []string       `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	Mentions []string       `protobuf:"bytes,5,rep,name=mentions,proto3" json:"mentions,omitempty"`
	Media    []*MediaUpload `protobuf:"bytes,6,rep,name=media,proto3" json:"media,omitempty"`
}

// Media upload information
type MediaUpload struct {
	FileData []byte `protobuf:"bytes,1,opt,name=file_data,json=fileData,proto3" json:"file_data,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileType string `protobuf:"bytes,3,opt,name=file_type,json=fileType,proto3" json:"file_type,omitempty"`
	PostId   string `protobuf:"bytes,4,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
}

// Response containing post data
type PostResponse struct {
	Post *Post `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
}

type GetPostRequest struct {
	PostId string `protobuf:"bytes,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
}

type GetUserPostsRequest struct {
	UserId  string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Page    int32  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PerPage int32  `protobuf:"varint,3,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
}

type GetUserPostsResponse struct {
	Posts []*Post `protobuf:"bytes,1,rep,name=posts,proto3" json:"posts,omitempty"` // Array of posts

}

type GetAllPostsRequest struct {
	Page    int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PerPage int32 `protobuf:"varint,2,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
}

type GetAllPostsResponse struct {
	Posts []*Post `protobuf:"bytes,1,rep,name=posts,proto3" json:"posts,omitempty"` // Array of posts

}

// Request to delete a post
type DeletePostRequest struct {
	PostId string `protobuf:"bytes,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // For authorization

}

// Response for delete operation
type DeletePostResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

type UpdatePostRequest struct {
	PostId  string `protobuf:"bytes,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // For authorization
	Caption string `protobuf:"bytes,3,opt,name=caption,proto3" json:"caption,omitempty"`
}

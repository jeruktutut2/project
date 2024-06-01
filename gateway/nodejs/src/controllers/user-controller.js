const register = async (req, res, next) => {
    console.log("register");
    res.status(200).json({data: "data", error: "error"})
}

const login = async (req, res, next) => {
    console.log("login");
    res.status(200).json({data: "data", error: "error"})
}

const logout = async (req, res, next) => {
    console.log(logout);
    res.status(200).json({data: "data", error: "error"})
}

export default {
    register,
    login,
    logout
}
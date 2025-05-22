import json
from typing import Any

from pydantic import BaseModel, model_validator


class ParseFromJsonMixin(BaseModel):
    @model_validator(mode="before")
    @classmethod
    def parse_from_json(cls, value: Any) -> Any:
        if isinstance(value, str):
            return json.loads(value)
        return value
